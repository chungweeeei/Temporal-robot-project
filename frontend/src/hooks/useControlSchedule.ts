import axios from "axios";
import { queryClient } from "@/utils/http";
import { useMutation } from "@tanstack/react-query";
import { type Schedule } from "@/types/schema";

async function pauseSchedule(scheduleId: string) {
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

async function unpauseSchedule(scheduleId: string) {
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


export function usePauseSchedule() {
  return useMutation({
    mutationFn: (scheduleId: string) => pauseSchedule(scheduleId),
    // Optimistic Update
    onMutate: async (scheduleId) => {
      // 取消正在進行的 refetch，避免覆蓋樂觀更新
      await queryClient.cancelQueries({ queryKey: ["schedules"] });

      // 保存之前的狀態以便回滾
      const previousSchedules = queryClient.getQueryData<Schedule[]>(["schedules"]);

      // 樂觀更新 cache
      queryClient.setQueryData<Schedule[]>(["schedules"], (old) =>
        old?.map((schedule) =>
          schedule.schedule_id === scheduleId
            ? { ...schedule, paused: true }
            : schedule
        )
      );

      return { previousSchedules };
    },
    onError: (_err, _scheduleId, context) => {
      // 發生錯誤時回滾
      if (context?.previousSchedules) {
        queryClient.setQueryData(["schedules"], context.previousSchedules);
      }
    },
    onSettled: async () => {
      await new Promise(resolve => setTimeout(resolve, 500));
      queryClient.invalidateQueries({ queryKey: ["schedules"] });
    },
  });
}

export function useUnpauseSchedule() {
  return useMutation({
    mutationFn: (scheduleId: string) => unpauseSchedule(scheduleId),
    onMutate: async (scheduleId) => {
      // 取消正在進行的 refetch，避免覆蓋樂觀更新
      await queryClient.cancelQueries({ queryKey: ["schedules"] });

      // 保存之前的狀態以便回滾
      const previousSchedules = queryClient.getQueryData<Schedule[]>(["schedules"]);

      // 樂觀更新 cache
      queryClient.setQueryData<Schedule[]>(["schedules"], (old) =>
        old?.map((schedule) =>
          schedule.schedule_id === scheduleId
            ? { ...schedule, paused: false }
            : schedule
        )
      );

      return { previousSchedules };
    },
    onError: (_err, _scheduleId, context) => {
      // 發生錯誤時回滾
      if (context?.previousSchedules) {
        queryClient.setQueryData(["schedules"], context.previousSchedules);
      }
    },
    onSettled: async () => {
      await new Promise(resolve => setTimeout(resolve, 500));
      queryClient.invalidateQueries({ queryKey: ["schedules"] });
    },
  });
}
