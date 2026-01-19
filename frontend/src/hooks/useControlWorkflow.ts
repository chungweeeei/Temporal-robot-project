import { useMutation } from "@tanstack/react-query";
import { pauseWorkflow, resumeWorkflow } from "../utils/http";


export const usePauseWorkflow = () => {
    return useMutation({
        mutationFn: (workflowId: string) => pauseWorkflow(workflowId),
    });
}

export const useResumeWorkflow = () => {
    return useMutation({
        mutationFn: (workflowId: string) => resumeWorkflow(workflowId),
    });
}