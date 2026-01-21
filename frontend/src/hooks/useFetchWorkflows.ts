import axios from "axios";
import { useQuery } from "@tanstack/react-query";
import { type WorkflowInfo } from "@/types/workflows";

async function fetchWorkflows(): Promise<WorkflowInfo[]> {

    const response = await axios.get<Promise<WorkflowInfo[]>>(
        "http://localhost:3000/api/v1/workflows",
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    )

    if (response.status !== 200){
        throw new Error(`Failed to fetch workflows: ${response.statusText}`);
    }

    return response.data;
}

async function fetchWorkflowById(workflowId: string){
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


export const useFetchWorkflow = () => {
    return useQuery({
        queryKey: ['workflows'],
        queryFn: fetchWorkflows,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });
}

export const useFetchWorkflowById = (workflowId: string, options?: { enabled?: boolean }) => {
    return useQuery({
        queryKey: ['workflow', workflowId],
        queryFn: () => fetchWorkflowById(workflowId),
        enabled: !!workflowId && (options?.enabled !== false),
        staleTime: 5 * 60 * 1000,
    });
}

