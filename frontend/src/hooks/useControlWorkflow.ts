import axios from "axios";
import { useMutation } from "@tanstack/react-query";

async function triggerWorkflow(workflowId: string){

    const response = await axios.post(
        `http://localhost:3000/api/v1/workflows/${workflowId}/trigger`,
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200){
        throw new Error(`Failed to trigger workflow: ${response.statusText}`);
    }

    return response.data;
}


async function pauseWorkflow(workflowId: string){

    const response = await axios.post(
        `http://localhost:3000/api/v1/workflows/${workflowId}/pause`,
        {},
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200){
        throw new Error(`Failed to pause workflow: ${response.statusText}`);
    }

    return response.data;
}

async function resumeWorkflow(workflowId: string){

    const response = await axios.post(
        `http://localhost:3000/api/v1/workflows/${workflowId}/resume`,
        {},
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200){
        throw new Error(`Failed to resume workflow: ${response.statusText}`);
    }

    return response.data;
}

export const useTriggerWorkflow = () => {
    return useMutation({
        mutationFn: (workflowId: string) => triggerWorkflow(workflowId),
    });
}


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