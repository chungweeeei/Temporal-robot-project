import { useMutation, useQueryClient } from "@tanstack/react-query";
import { deleteWorkflow } from "@/utils/http";

export const useDeleteWorkflow = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (workflowId: string) => deleteWorkflow(workflowId),
        onSuccess: () => {
            // 刪除成功後，重新抓取 workflow 列表
            queryClient.invalidateQueries({ queryKey: ['workflows'] });
        }
    });
};
