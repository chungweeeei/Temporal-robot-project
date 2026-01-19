import axios from "axios";
import { QueryClient } from "@tanstack/react-query";
import type { CreateSchedulePayload } from "../types/schema";

// 建立一個共用的 Query Client 實例
export const queryClient = new QueryClient();

export async function createSchedule(data: CreateSchedulePayload) {
    const response = await axios.post(
        "http://localhost:3000/api/v1/schedules",
        data,
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200 && response.status !== 201) {
        throw new Error(`Failed to create schedule: ${response.statusText}`);
    }

    return response.data;
}

export async function fetchSchedules() {
    const response = await axios.get(
        "http://localhost:3000/api/v1/schedules",
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200) {
        throw new Error(`Failed to fetch schedules: ${response.statusText}`);
    }

    return response.data;
}

export async function pauseSchedule(scheduleId: string) {
    const response = await axios.post(
        `http://localhost:3000/api/v1/schedules/${scheduleId}/pause`,
        {},
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200) {
        throw new Error(`Failed to pause schedule: ${response.statusText}`);
    }

    return response.data;
}

export async function unpauseSchedule(scheduleId: string) {
    const response = await axios.post(
        `http://localhost:3000/api/v1/schedules/${scheduleId}/resume`,
        {},
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200) {
        throw new Error(`Failed to unpause schedule: ${response.statusText}`);
    }

    return response.data;
}