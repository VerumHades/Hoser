import { useEffect, useState } from "react";
import { DeveloperListingAPI, type ListingGithubSetup } from "../../../backend/repositories/developer_listing";
import toast from "react-hot-toast";
import Section from "../../../templates/components/Section";
import { Button } from "../../../templates/components/Button";

export function GithubSetupForm({ listingId }: { listingId: string }) {
    const [setup, setSetup] = useState<ListingGithubSetup | null>(null);
    const [draft, setDraft] = useState({ repoUrl: "", accessToken: "" });
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        DeveloperListingAPI.githubSetup.get(listingId).then(data => {
            setSetup(data);
            if (data) setDraft({ repoUrl: data.repositoryURL, accessToken: data.accessToken });
            setLoading(false);
        });
    }, [listingId]);

    const handleSave = async () => {
        setLoading(true);
        const res = await DeveloperListingAPI.githubSetup.attachOrUpdate(listingId, draft.repoUrl, draft.accessToken);
        setSetup(res);
        setLoading(false);
        toast.success("GitHub attached");
    };

    return (
        <Section title="GitHub Setup" description="Automated deployment repository.">
            <div className="flex flex-col gap-2">
                <input className="..." value={draft.repoUrl} onChange={e => setDraft({...draft, repoUrl: e.target.value})} placeholder="Repo URL" />
                <input className="..." value={draft.accessToken} onChange={e => setDraft({...draft, accessToken: e.target.value})} placeholder="Token" />
                <div className="flex gap-2">
                    <Button variant="success" onClick={handleSave} disabled={loading}>{setup ? "Update" : "Attach"}</Button>
                    {setup && <Button variant="danger" onClick={async () => { 
                        await DeveloperListingAPI.githubSetup.remove(listingId); 
                        setSetup(null); setDraft({repoUrl: "", accessToken: ""});
                    }}>Remove</Button>}
                </div>
            </div>
        </Section>
    );
}