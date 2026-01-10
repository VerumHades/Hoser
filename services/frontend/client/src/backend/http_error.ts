/**
 * Represents an HTTP error with a status code.
 */
export class HttpError extends Error {
    public readonly statusCode: number;

    /**
     * Creates a new HttpError.
     * @param statusCode - HTTP status code (e.g., 404, 500)
     * @param message - Error message
     */
    constructor(statusCode: number, message: string) {
        super(message);
        this.statusCode = statusCode;
        this.name = "HttpError";
        Object.setPrototypeOf(this, HttpError.prototype); // Required for instanceof to work
    }
}
