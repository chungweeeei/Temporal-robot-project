import axios from "axios";
import { useQuery } from "@tanstack/react-query";

export async function fetchWorkflowById(workflowId: string){
    const response = await axios.get(
        `http://localhost:3000/api/v1/workflows/${workflowId}`,
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    )

    if (response.status !== 200){
        throw new Error(`Failed to fetch workflow: ${response.statusText}`);
    }

    return response.data;
}

export const useFetchWorkflowById = (workflowId: string, options?: { enabled?: boolean }) => {
    return useQuery({
        queryKey: ['workflow', workflowId],
        queryFn: () => fetchWorkflowById(workflowId),
        enabled: !!workflowId && (options?.enabled !== false),
        staleTime: 5 * 60 * 1000,
    });
}