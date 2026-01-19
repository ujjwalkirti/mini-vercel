import { ClickHouseClient } from "@clickhouse/client";
import { v4 as uuidv4 } from "uuid";

export interface ClickHouseInsertValues {
    deployment_id: string;
    log: string;
}

export default class ClickHouseService {
    private clickHouseClient: ClickHouseClient;

    constructor(clickHouseClient: ClickHouseClient) {
        this.clickHouseClient = clickHouseClient;
    }

    async insertLog(
        tableName: string,
        values: ClickHouseInsertValues
    ): Promise<{ query_id: string }> {
        const { deployment_id, log } = values;

        if (!deployment_id || !log) return { query_id: "" };

        const { query_id } = await this.clickHouseClient.insert({
            table: tableName,
            values: [
                {
                    event_id: uuidv4(),
                    deployment_id,
                    log
                }
            ],
            format: "JSONEachRow"
        });

        return { query_id };
    }
}
