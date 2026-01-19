import axios from "axios";
import { useMutation } from "@tanstack/react-query";
import { type NodeInfo } from "@/types/workflows";

type TriggerWorkflowPayload = {
    workflow_id: string;
    workflow_name: string;
    nodes: Record<string, NodeInfo>;
}

async function triggerWorkflow(payload: TriggerWorkflowPayload){

    const response = await axios.post(
        `http://localhost:3000/api/v1/workflows/${payload.workflow_id}/trigger`,
        {
            ...payload,
            root_node_id: "start",
        },
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
        mutationFn: (payload: TriggerWorkflowPayload) => triggerWorkflow(payload),
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