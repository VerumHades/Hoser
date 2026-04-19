package external

import (
	"common/pkg/domain/entities/listing"
	"common/pkg/shared"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

type MeiliListingSearchIndex struct {
	client meilisearch.ServiceManager
	index  meilisearch.IndexManager
}

// NewMeiliListingSearchIndex creates a new search indexer and ensures proper configuration.
func NewMeiliListingSearchIndex(host string, apiKey string) *MeiliListingSearchIndex {
	client := meilisearch.New(host, meilisearch.WithAPIKey(apiKey))
	setupIndex(client, "listings")

	return &MeiliListingSearchIndex{
		client: client,
		index:  client.Index("listings"),
	}
}

func setupIndex(client meilisearch.ServiceManager, indexUID string) {
	_, err := client.GetIndex(indexUID)
	if err != nil {
		createIndex(client, indexUID)
	}
	configureIndexSettings(client, indexUID)
}

func createIndex(client meilisearch.ServiceManager, indexUID string) {
	config := &meilisearch.IndexConfig{Uid: indexUID, PrimaryKey: "id"}
	_, _ = client.CreateIndex(config)
}

func configureIndexSettings(client meilisearch.ServiceManager, indexUID string) {
	settings := &meilisearch.Settings{
		SearchableAttributes: []string{"title", "description", "authorID"},
		FilterableAttributes: []string{"accessMode", "cpuCount", "ramBytes", "price", "createdAtTimestamp"},
		SortableAttributes:   []string{"createdAtTimestamp", "price"},
	}
	_, _ = client.Index(indexUID).UpdateSettings(settings)
}

// Index adds a new listing to the search engine.
func (indexer *MeiliListingSearchIndex) Index(ctx context.Context, listingToSave *listing.Listing) error {
	document := mapListingToDocument(listingToSave)
	_, err := indexer.index.AddDocuments(document, nil)
	return err
}

/**
 * Update modifies an existing listing in the search index using a partial update event.
 * It maps the event into a document containing only the non-nil fields.
 */
func (indexer *MeiliListingSearchIndex) Update(
	ctx context.Context,
	listingUpdate *listing.ListingUpdateEvent,
) error {
	updateDocument := mapUpdateEventToDocument(listingUpdate)

	// Meilisearch's UpdateDocuments performs a partial update (merge)
	// when provided with the primary key.
	_, err := indexer.index.UpdateDocuments(updateDocument, nil)
	return err
}

/**
 * Maps a ListingUpdateEvent to a map of changes for the search engine.
 * Only includes fields that are not nil to support partial updates.
 */
func mapUpdateEventToDocument(event *listing.ListingUpdateEvent) map[string]interface{} {
	document := make(map[string]interface{})

	// The ID is mandatory as the primary key for the update operation.
	document["id"] = string(event.ID)

	applyTextUpdates(document, event)
	applyNumericUpdates(document, event)
	applyHardwareUpdates(document, event)
	applyCollectionUpdates(document, event)

	return document
}

func applyTextUpdates(document map[string]interface{}, event *listing.ListingUpdateEvent) {
	if event.Title != nil {
		document["title"] = *event.Title
	}
	if event.Description != nil {
		document["description"] = *event.Description
	}
	if event.DocumentationMarkdown != nil {
		document["documentationMarkdown"] = *event.DocumentationMarkdown
	}
}

func applyNumericUpdates(document map[string]interface{}, event *listing.ListingUpdateEvent) {
	if event.AccessMode != nil {
		document["accessMode"] = int(*event.AccessMode)
	}
	if event.PriceInMinorUnits != nil {
		document["price"] = *event.PriceInMinorUnits
	}
}

func applyHardwareUpdates(document map[string]interface{}, event *listing.ListingUpdateEvent) {
	if event.HardwareSpecification != nil {
		document["cpuCount"] = event.HardwareSpecification.CPUCount
		document["ramBytes"] = event.HardwareSpecification.RAMBytes
		document["diskBytes"] = event.HardwareSpecification.DiskBytes
	}
}

func applyCollectionUpdates(document map[string]interface{}, event *listing.ListingUpdateEvent) {
	if event.ScreenshotKeys != nil {
		document["screenshotKeys"] = *event.ScreenshotKeys
	}
}

// Remove deletes a listing from the search index.
func (indexer *MeiliListingSearchIndex) Remove(ctx context.Context, listingID shared.ListingID) error {
	_, err := indexer.index.DeleteDocument(string(listingID), nil)
	return err
}

// SearchNextBatch retrieves a paginated set of listings or returns ErrNoContentFound.
func (indexer *MeiliListingSearchIndex) SearchNextBatch(
	ctx context.Context,
	query listing.SearchQuery,
	batchRequest shared.BatchRequest[listing.ListingSearchCursor],
) ([]*listing.Listing, listing.ListingSearchCursor, error) {
	searchRequest := indexer.buildSearchRequest(query, batchRequest)

	response, err := indexer.index.Search(query.Text, searchRequest)
	if err != nil {
		return nil, listing.ListingSearchCursor{}, err
	}

	if len(response.Hits) == 0 {
		return nil, listing.ListingSearchCursor{}, shared.ErrNoContentFound
	}

	return indexer.processSearchResults(response)
}

func (indexer *MeiliListingSearchIndex) buildSearchRequest(
	query listing.SearchQuery,
	request shared.BatchRequest[listing.ListingSearchCursor],
) *meilisearch.SearchRequest {
	filters := indexer.assembleFilterList(query, request.Cursor)

	return &meilisearch.SearchRequest{
		Limit:  int64(request.MaxBatchSize),
		Filter: strings.Join(filters, " AND "),
		Sort:   []string{"createdAtTimestamp:asc"},
	}
}

func (indexer *MeiliListingSearchIndex) assembleFilterList(
	query listing.SearchQuery,
	cursor listing.ListingSearchCursor,
) []string {
	filters := []string{
		fmt.Sprintf("createdAtTimestamp > %d", cursor.LastCreatedAt.Unix()),
	}

	filters = append(filters, appendRangeFilter("cpuCount", query.CPU)...)
	filters = append(filters, appendRangeFilter("ramBytes", query.RAMBytes)...)
	filters = append(filters, appendRangeFilter("price", query.Price)...)
	filters = append(filters, appendDateFilter("createdAtTimestamp", query.CreatedAt)...)

	if query.AccessMode != nil {
		filters = append(filters, fmt.Sprintf("accessMode = %d", *query.AccessMode))
	}

	return filters
}

func appendRangeFilter(field string, numericRange listing.NumericRange) []string {
	var filters []string
	if numericRange.Min != nil {
		filters = append(filters, fmt.Sprintf("%s >= %d", field, *numericRange.Min))
	}
	if numericRange.Max != nil {
		filters = append(filters, fmt.Sprintf("%s <= %d", field, *numericRange.Max))
	}
	return filters
}

func appendDateFilter(field string, dateRange listing.DateRange) []string {
	var filters []string
	if dateRange.From != nil {
		filters = append(filters, fmt.Sprintf("%s >= %d", field, dateRange.From.Unix()))
	}
	if dateRange.To != nil {
		filters = append(filters, fmt.Sprintf("%s <= %d", field, dateRange.To.Unix()))
	}
	return filters
}

func (indexer *MeiliListingSearchIndex) processSearchResults(
	response *meilisearch.SearchResponse,
) ([]*listing.Listing, listing.ListingSearchCursor, error) {
	listings := make([]*listing.Listing, 0, len(response.Hits))

	for _, hit := range response.Hits {
		listings = append(listings, decodeHitToListing(hit))
	}

	return listings, calculateNextCursor(listings), nil
}

func decodeHitToListing(hit interface{}) *listing.Listing {
	rawJSON, _ := json.Marshal(hit)
	var document map[string]interface{}
	_ = json.Unmarshal(rawJSON, &document)

	return mapDocumentToListing(document)
}

func calculateNextCursor(listings []*listing.Listing) listing.ListingSearchCursor {
	var cursor listing.ListingSearchCursor

	if len(listings) > 0 {
		cursor.LastCreatedAt = listings[len(listings)-1].CreatedAt()
	}

	return cursor
}

func mapListingToDocument(listingEntity *listing.Listing) map[string]interface{} {
	hardware := listingEntity.HardwareSpecification()
	return map[string]interface{}{
		"id":                 string(listingEntity.ID()),
		"authorID":           string(listingEntity.AuthorID()),
		"title":              listingEntity.Title(),
		"description":        listingEntity.Description(),
		"accessMode":         int(listingEntity.AccessMode()),
		"cpuCount":           hardware.CPUCount,
		"ramBytes":           hardware.RAMBytes,
		"diskBytes":          hardware.DiskBytes,
		"price":              listingEntity.PriceInMinorUnits(),
		"createdAtTimestamp": listingEntity.CreatedAt().Unix(),
		"screenshotKeys":     listingEntity.ScreenshotKeys(),
	}
}

func mapDocumentToListing(document map[string]interface{}) *listing.Listing {
	specification := &shared.HardwareSpecification{
		CPUCount:  int64(document["cpuCount"].(float64)),
		RAMBytes:  int64(document["ramBytes"].(float64)),
		DiskBytes: int64(document["diskBytes"].(float64)),
	}

	listingEntity, _, _ := listing.NewListingWithID(
		shared.ListingID(document["id"].(string)),
		shared.UserID(document["authorID"].(string)),
		document["title"].(string),
		document["description"].(string),
		listing.ListingAccessMode(document["accessMode"].(float64)),
		specification,
		int64(document["price"].(float64)),
		time.Unix(int64(document["createdAtTimestamp"].(float64)), 0),
		hydrateScreenshots(document["screenshotKeys"]),
		"",
	)

	//fmt.Println(err)

	return listingEntity
}

func hydrateScreenshots(rawIDs interface{}) []shared.ListingScreenshotID {
	items, ok := rawIDs.([]interface{})
	if !ok {
		return nil
	}

	screenshotIDs := make([]shared.ListingScreenshotID, len(items))
	for index, value := range items {
		screenshotIDs[index] = shared.ListingScreenshotID(value.(string))
	}
	return screenshotIDs
}
