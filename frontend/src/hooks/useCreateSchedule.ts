import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createSchedule } from "@/utils/http";
import { type CreateSchedulePayload } from "@/types/schema";

export function useCreateSchedule() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateSchedulePayload) => {
        return createSchedule(data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["schedules"] });
    },
  });
}
