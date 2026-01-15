import { useMutation } from "@tanstack/react-query";
import { triggerWorkflow } from "../utils/http";
import type { WorkflowPayload } from "../types/schema";

export const useTriggerWorkflow = () => {
    return useMutation({
        mutationFn: (data: WorkflowPayload) => triggerWorkflow(data),
        onSuccess: () => {
            console.log("Workflow triggered successfully");
        }
    });
}