export type ActivityType = "Standup" | "Standdown" | "Sitdown" | "Move" | "Sleep" | "Start" | "End" | "TTS" | "Head";
export type WorkflowStatusDef = "Idle" | "Running" | "Completed" | "Failed" | "Paused" | "Cancelled" | "Terminated";

export interface BaseParams {
  [key: string]: string | number | boolean;
}

export interface MoveParams extends BaseParams {
  x: number;
  y: number;
  orientation: number;
}

export interface SleepParams extends BaseParams {
  duration: number; // milliseconds
}

export interface HeadParams extends BaseParams {
  angle: number; // degrees
}

export interface TTSParams extends BaseParams {
  text: string; // text to speak
}

export interface NodeInfo {
    id: string;
    type: ActivityType | "Start" | "End";
    params: BaseParams;
    transitions: {
        next?: string;
        failure?: string;
    };
}

export interface WorkflowInfo {
    workflow_id: string;
    workflow_name: string;
    root_node_id: string;
    nodes: Record<string, NodeInfo>;
    created_at: number;
    updated_at: number;
}

export interface WorkflowStatus {
    status: WorkflowStatusDef;
    current_step: string;
}


