import axios from "axios";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export async function createWorkflow(workflowName: string){

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
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (workflowName: string) => createWorkflow(workflowName),
        onSuccess: () => {
            // 建立成功後，重新抓取 workflow 列表
            queryClient.invalidateQueries({ queryKey: ['workflows'] });
        }
    });
};
