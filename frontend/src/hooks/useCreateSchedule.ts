import axios from "axios";
import { queryClient } from "@/utils/http";
import { useMutation } from "@tanstack/react-query";
import { type CreateSchedulePayload } from "@/types/schema";

async function createSchedule(data: CreateSchedulePayload) {
    const response = await axios.post(
        "http://localhost:3000/api/v1/schedules",
        {
          ...data,
          timezone: "Asia/Taipei",
        },
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

export function useCreateSchedule() {

  return useMutation({
    mutationFn: (data: CreateSchedulePayload) => {
        return createSchedule(data);
    },
    onSuccess: async() => {
      await new Promise(resolve => setTimeout(resolve, 300));
      queryClient.invalidateQueries({ queryKey: ["schedules"] });
    },
  });
}
