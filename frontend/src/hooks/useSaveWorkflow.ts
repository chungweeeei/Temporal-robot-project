import { useMutation } from "@tanstack/react-query";
import { saveWorkflow } from "../utils/http";
import type { WorkflowPayload } from "../types/schema";
import { queryClient } from "../utils/http";

export const useSaveWorkflow = () => {
    return useMutation({
        mutationFn: (data: WorkflowPayload) => saveWorkflow(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['workflows'] });
        }
    });
}