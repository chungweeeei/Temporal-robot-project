import axios from "axios";
import { queryClient } from "@/utils/http";
import { useMutation } from "@tanstack/react-query";

async function createWorkflow(workflowName: string){

    const response = await axios.post(
        "http://localhost:3000/api/v1/workflows",
        {
            workflow_name: workflowName,
            nodes: {}
        },
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200 && response.status !== 201){
        throw new Error(`Failed to create workflow: ${response.statusText}`);
    }

    return response.data;
}

export const useCreateWorkflow = () => {
    return useMutation({
        mutationFn: (workflowName: string) => createWorkflow(workflowName),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['workflows'] });
        }
    });
};
