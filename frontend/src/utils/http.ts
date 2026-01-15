import axios from "axios";
import { QueryClient } from "@tanstack/react-query";
import type { WorkflowPayload } from "../types/schema";

// 建立一個共用的 Query Client 實例
export const queryClient = new QueryClient();

export async function fetchWorkflows(){

    const response = await axios.get(
        "http://localhost:3000/api/v1/workflows",
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    )

    if (response.status !== 200){
        throw new Error(`Failed to fetch workflows: ${response.statusText}`);
    }

    return response.data;
}

export async function fetchWorkflowById(workflowId: string){

    const response = await axios.get(
        `http://localhost:3000/api/v1/workflows/${workflowId}`,
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    )

    if (response.status !== 200){
        throw new Error(`Failed to fetch workflow: ${response.statusText}`);
    }

    return response.data;
}

export async function fetchWorkflowStatus(workflowId: string){

    const response = await axios.get(
        `http://localhost:3000/api/v1/workflows/${workflowId}/status`,
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    )

    if (response.status !== 200){
        throw new Error(`Failed to fetch workflow status: ${response.statusText}`);
    }
    
    return response.data;
}


export async function saveWorkflow(data: WorkflowPayload){

    const response = await axios.post(
        "http://localhost:3000/api/v1/workflows",
        data,
        {
            headers:{
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200){
        throw new Error(`Failed to save workflow: ${response.statusText}`);
    }

    return response.data;
}

export async function triggerWorkflow(data: WorkflowPayload){

    const response = await axios.post(
        `http://localhost:3000/api/v1/workflows/${data.workflow_id}/trigger`,
        data,
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

export async function pauseWorkflow(workflowId: string){

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

export async function resumeWorkflow(workflowId: string){

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