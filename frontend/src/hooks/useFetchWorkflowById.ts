import { useQuery } from "@tanstack/react-query";
import { fetchWorkflowById } from "../utils/http"

export const useFetchWorkflowById = (workflowId: string, options?: { enabled?: boolean }) => {
    return useQuery({
        queryKey: ['workflow', workflowId],
        queryFn: () => fetchWorkflowById(workflowId),
        enabled: !!workflowId && (options?.enabled !== false),
        staleTime: 5 * 60 * 1000,
    });
}