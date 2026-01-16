import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createWorkflow } from "@/utils/http";

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
