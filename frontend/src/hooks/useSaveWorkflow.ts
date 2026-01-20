import axios from "axios";
import { useMutation } from "@tanstack/react-query";
import type { NodeInfo } from "../types/workflows";
import { queryClient } from "../utils/http";

export type SaveWorkflowPayload = {
    workflow_id: string;
    workflow_name: string;
    nodes: Record<string, NodeInfo>;
}

async function saveWorkflow(payload: SaveWorkflowPayload){
    const response = await axios.post(
        "http://localhost:3000/api/v1/workflows",
        {
            ...payload,
            root_node_id: "start",
        },
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

export const useSaveWorkflow = () => {
    return useMutation({
        mutationFn: (payload: SaveWorkflowPayload) => saveWorkflow(payload),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['workflows'] });
        }
    });
}