
// 這是 React Flow Node 的 data 結構
// export interface FlowNodeData extends Record<string, unknown> {
//   label: string;
//   activityType: ActivityType;
//   params: MoveParams | SleepParams | BaseParams | HeadParams | TTSParams;
// }


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
