import { useQuery } from "@tanstack/react-query";
import { fetchWorkflowById } from "../utils/http"

export const useFetchWorkflowById = (workflowId: string) => {
    return useQuery({
        queryKey: ['workflow', workflowId],
        queryFn: () => fetchWorkflowById(workflowId),
        enabled: !!workflowId, // 只有當 ID 存在時才執行 Query
        staleTime: 5 * 60 * 1000, // 5 minutes
    });
}