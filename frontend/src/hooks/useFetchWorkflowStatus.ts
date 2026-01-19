import axios from "axios";
import { useQuery } from "@tanstack/react-query";
import type { WorkflowStatus } from "@/types/workflows";


export async function fetchWorkflowStatus(workflowId: string): Promise<WorkflowStatus> {

    const response = await axios.get<WorkflowStatus>(
        `http://localhost:3000/api/v1/workflows/${workflowId}/status`,
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    )

    if (response.status !== 200){
        throw new Error(`Failed to fetch workflow status: ${response.statusText}`);
    }
    
    return response.data;
}


export const useFetchWorkflowStatus = (workflowId: string, isRunning: boolean) => {
    return useQuery({
        queryKey: ['workflows', workflowId, "status"],
        queryFn: () => fetchWorkflowStatus(workflowId),
        enabled: !!workflowId && isRunning,
        refetchInterval: isRunning ? 2000 : false,
        staleTime: 0
    });
}
