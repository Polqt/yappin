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
	try {
		const errorData = await response.json();
		const errorMessage = errorData.error || errorData.message || 'Request failed';
		throw new Error(errorMessage);
	} catch (e) {
		if (e instanceof Error && e.message !== 'Request failed') {
			throw e;
		}
		throw new Error(`Request failed with status ${response.status}`);
	}
}
