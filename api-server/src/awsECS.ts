import { AssignPublicIp, ECSClient, LaunchType, RunTaskCommand } from "@aws-sdk/client-ecs";
import { EnvVar } from "./azureACIREST.js";

interface ECSTask {
    cluster: string,
    taskDefinition: string,
    launchType: string,
    count: number,
    subnets: string[],
    securityGroups: string[],
    assignPublicIp: string,
    image: string,
    envVars: EnvVar[]
}

class AWSECSService {
    private ecsClient: ECSClient;
    constructor(ecsClient: ECSClient) {
        this.ecsClient = ecsClient
    }

    async runTask({ cluster, taskDefinition, launchType, count, subnets, securityGroups, assignPublicIp, image, envVars }: ECSTask) {
        const params = {
            cluster: cluster,
            taskDefinition: taskDefinition,
            launchType: launchType as LaunchType,
            count: count,
            networkConfiguration: {
                awsvpcConfiguration: {
                    subnets: subnets,
                    securityGroups: securityGroups,
                    assignPublicIp: assignPublicIp as AssignPublicIp
                }
            },
            overrides: {
                containerOverrides: [
                    {
                        name: image,
                        environment: envVars
                    }
                ]
            }
        }
        const command = new RunTaskCommand(params);

        await this.ecsClient.send(command);
    }
}


export default AWSECSService
