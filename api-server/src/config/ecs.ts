export const ecsConfig = {
    cluster: process.env.ECS_CLUSTER_NAME ?? "",
    taskDefinition: process.env.ECS_TASK_DEFINITION ?? "",
    imageName: process.env.ECS_IMAGE_NAME ?? "",
    subnets: process.env.ECS_SUBNETS?.split(",").map((subnet) => subnet.trim()) ?? [],
    securityGroups: process.env.ECS_SECURITY_GROUPS?.split(",").map((sg) => sg.trim()) ?? [],
    launchType: process.env.ECS_LAUNCH_TYPE ?? "FARGATE",
    count: parseInt(process.env.ECS_COUNT ?? "1"),
    assignPublicIp: process.env.ECS_ASSIGN_PUBLIC_IP ?? "ENABLED",
}
