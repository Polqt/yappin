const isDev = import.meta.env.DEV;

export const logger = {
	/**
	 * Log general information (disabled in production)
	 */
	log: (...args: unknown[]): void => {
		if (isDev) console.log(...args);
	},

	/**
	 * Log debug information (disabled in production)
	 */
	debug: (...args: unknown[]): void => {
		if (isDev) console.debug(...args);
	},

	/**
	 * Log warnings (disabled in production)
	 */
	warn: (...args: unknown[]): void => {
		if (isDev) console.warn(...args);
	},

	/**
	 * Log errors (always enabled)
	 */
	error: (...args: unknown[]): void => {
		console.error(...args);
	},

	/**
	 * Log with a specific prefix for easier filtering
	 */
	prefix: (prefix: string) => ({
		log: (...args: unknown[]) => logger.log(`[${prefix}]`, ...args),
		debug: (...args: unknown[]) => logger.debug(`[${prefix}]`, ...args),
		warn: (...args: unknown[]) => logger.warn(`[${prefix}]`, ...args),
		error: (...args: unknown[]) => logger.error(`[${prefix}]`, ...args)
	})
};

// Pre-configured loggers for common modules
export const wsLogger = logger.prefix('WebSocket');
export const authLogger = logger.prefix('Auth');
export const roomLogger = logger.prefix('Room');
