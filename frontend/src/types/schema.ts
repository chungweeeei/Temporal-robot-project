export type ActivityType = 'Standup' | 'Sitdown' | 'Move' | 'Sleep' | 'Start' | 'End' | 'TTS' | 'Head';

export type WorkflowStatus = 'idle' | 'running' | 'paused' | 'completed' | 'failed';

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

// 這是 React Flow Node 的 data 結構
export interface FlowNodeData extends Record<string, unknown> {
  label: string;
  activityType: ActivityType;
  params: MoveParams | SleepParams | BaseParams | HeadParams | TTSParams;
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

// Dashboard 用的 Workflow 摘要
export interface WorkflowSummary {
  workflow_id: string;
  workflow_name: string;
  status?: WorkflowStatus;
  current_step?: string;
  updated_at?: string;
}

// 排程功能 (預留)
export interface WorkflowSchedule {
  enabled: boolean;
  cron_expression?: string;  // e.g., "0 8 * * *" for daily 8 AM
  next_run?: string;         // ISO timestamp
  last_run?: string;         // ISO timestamp
}

export interface CreateSchedulePayload {
  schedule_id: string;
  workflow_id: string;
  cron_expr: string;
  timezone?: string;
}

export interface ScheduleRange {
  start: number;
  end: number;
  step: number;
}

export interface ScheduleCalendarSpec {
  second?: ScheduleRange[];
  minute?: ScheduleRange[];
  hour?: ScheduleRange[];
  day_of_month?: ScheduleRange[];
  month?: ScheduleRange[];
  day_of_week?: ScheduleRange[];
  year?: ScheduleRange[];
  comment?: string;
}

export interface ScheduleSpec {
  calendars?: ScheduleCalendarSpec[];
  cron_expressions?: string[];
  timezone_name?: string;
}

export interface Schedule {
  schedule_id: string;
  paused: boolean;
  spec?: ScheduleSpec;
  
  recent_run?: string;
  upcoming_run?: string;

  // Legacy/Flattened fields if still used elsewhere
  workflow_id?: string;
  cron_expr?: string;
  timezone?: string;
}
