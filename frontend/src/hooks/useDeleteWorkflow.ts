import { useMutation, useQueryClient } from "@tanstack/react-query";
import { deleteWorkflow } from "@/utils/http";

export const useDeleteWorkflow = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (workflowId: string) => deleteWorkflow(workflowId),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['workflows'] });
        }
    });
};
