import axios from "axios";
import { useQuery } from "@tanstack/react-query";
import type { Schedule } from "@/types/schema";

async function fetchSchedules(): Promise<Schedule[]> {
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

async function fetchSchedulesById(scheduleId: string): Promise<Schedule> {
    const response = await axios.get<Promise<Schedule>>(
        `http://localhost:3000/api/v1/schedules/${scheduleId}`,
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );
    
    if (response.status !== 200) {
        throw new Error(`Failed to fetch schedule: ${response.statusText}`);
    }

    return response.data;
}

export function useFetchSchedules() {
  return useQuery<Schedule[]>({
    queryKey: ["schedules"],
    queryFn: fetchSchedules,
  });
}

export function useFetchSchedulesById(scheduleId: string) {
  return useQuery<Schedule>({
    queryKey: ["schedules", scheduleId],
    queryFn: () => fetchSchedulesById(scheduleId),
  });
}

