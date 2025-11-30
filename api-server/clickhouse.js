const { ClickHouseClient } = require("@clickhouse/client");

class ClickHouseService {
    /**
     * Constructs a new CLickHouseService instance.
     * @param {ClickHouseClient} clickHouseClient - An instance of the ClickHouse class.
     */
    constructor(clickHouseClient) {
        this.clickHouseClient = clickHouseClient;
    }

    async insertData(table_name, values) {
        const { deployment_id, log } = values;
        if (!deployment_id || !log) return;
        const { query_id } = await this.clickHouseClient.insert({
            table: table_name,
            values: [{ event_id: uuidv4(), deployment_id: deployment_id, log }],
            format: 'JSONEachRow'
        })
        return query_id
    }
}


module.exports = ClickHouseService
