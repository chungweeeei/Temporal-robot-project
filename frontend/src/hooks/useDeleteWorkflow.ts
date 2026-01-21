import axios from "axios";
import { queryClient } from "@/utils/http";
import { useMutation } from "@tanstack/react-query";

async function deleteWorkflow(workflowId: string){

    const response = await axios.delete(
        `http://localhost:3000/api/v1/workflows/${workflowId}`,
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200 && response.status !== 204){
        throw new Error(`Failed to delete workflow: ${response.statusText}`);
    }

    return response.data;
}

export const useDeleteWorkflow = () => {
    return useMutation({
        mutationFn: (workflowId: string) => deleteWorkflow(workflowId),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['workflows'] });
        }
    });
};
