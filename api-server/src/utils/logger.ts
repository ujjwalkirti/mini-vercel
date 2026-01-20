type LogLevel = 'debug' | 'info' | 'warn' | 'error';

interface LogEntry {
    timestamp: string;
    level: LogLevel;
    message: string;
    context?: Record<string, unknown>;
}

class Logger {
    private formatLog(entry: LogEntry): string {
        const contextStr = entry.context ? ` ${JSON.stringify(entry.context)}` : '';
        return `[${entry.timestamp}] [${entry.level.toUpperCase()}] ${entry.message}${contextStr}`;
    }

    private log(level: LogLevel, message: string, context?: Record<string, unknown>): void {
        const entry: LogEntry = {
            timestamp: new Date().toISOString(),
            level,
            message,
            context
        };

        const formattedLog = this.formatLog(entry);

        switch (level) {
            case 'debug':
                if (process.env.NODE_ENV === 'development') {
                    console.debug(formattedLog);
                }
                break;
            case 'info':
                console.info(formattedLog);
                break;
            case 'warn':
                console.warn(formattedLog);
                break;
            case 'error':
                console.error(formattedLog);
                break;
        }
    }

    debug(message: string, context?: Record<string, unknown>): void {
        this.log('debug', message, context);
    }

    info(message: string, context?: Record<string, unknown>): void {
        this.log('info', message, context);
    }

    warn(message: string, context?: Record<string, unknown>): void {
        this.log('warn', message, context);
    }

    error(message: string, error?: Error | unknown, context?: Record<string, unknown>): void {
        const errorContext = error instanceof Error
            ? { errorMessage: error.message, stack: error.stack, ...context }
            : { error, ...context };
        this.log('error', message, errorContext);
    }
}

export const logger = new Logger();
