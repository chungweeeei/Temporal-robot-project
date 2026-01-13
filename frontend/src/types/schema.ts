export type ActivityType = 'Standup' | 'Sitdown' | 'Move' | 'Sleep' | 'Start' | 'End';

export interface RetryPolicy {
  max_attempts: number;
  initial_interval: number; // ms
  backoff_coefficient?: number;
  maximum_interval?: number; // ms
}

export interface BaseParams {
  [key: string]: string | number | boolean;
}

export interface MoveParams extends BaseParams {
  x: number;
  y: number;
}

export interface SleepParams extends BaseParams {
  duration: number; // milliseconds
}

export interface ConditionParams extends BaseParams {
  expression: string; // e.g. "x > 5"
}

// 這是 React Flow Node 的 data 結構
export interface FlowNodeData extends Record<string, unknown> {
  label: string;
  activityType: ActivityType;
  params: MoveParams | SleepParams | ConditionParams | BaseParams;
  retry_policy?: RetryPolicy;
}

export interface WorkflowTransitions {
  next?: string;
  failure?: string;
  true?: string;
  false?: string;
}

// 這是要送給後端的單一節點結構
export interface WorkflowNode {
  id: string;
  type: ActivityType;
  params: BaseParams;
  retry_policy?: RetryPolicy;
  transitions: WorkflowTransitions;
}

// 最終 Payload
export interface WorkflowPayload {
  workflow_id?: string;
  workflow_name?: string;
  root_node_id?: string;
  nodes: Record<string, WorkflowNode>;
}
