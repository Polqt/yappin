export function getErrorMessage(error: unknown): string {
	if (error instanceof Error) {
		return error.message;
	}
	if (typeof error === 'string') {
		return error;
	}
	return 'An unexpected error occurred';
}

export async function handleApiError(response: Response): Promise<never> {
	const errorText = await response.text().catch(() => 'Request failed');
	throw new Error(errorText || `Request failed with status ${response.status}`);
}
