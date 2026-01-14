import { useQuery } from "@tanstack/react-query";
import { fetchWorkflowStatus } from "../utils/http";

export const useFetchWorkflowStatus = (workflowId: string, isRunning: boolean) => {
    return useQuery({
        queryKey: ['workflows', workflowId, "status"],
        queryFn: () => fetchWorkflowStatus(workflowId),
        enabled: !!workflowId && isRunning,
        refetchInterval: isRunning ? 2000 : false,
        staleTime: 0
    });
}
