import axios from "axios";
import { queryClient } from "@/utils/http";
import { useMutation } from "@tanstack/react-query";

async function deleteSchedule(scheduleId: string){

    const response = await axios.delete(
        `http://localhost:3000/api/v1/schedules/${scheduleId}`,
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200 && response.status !== 204){
        throw new Error(`Failed to delete schedule: ${response.statusText}`);
    }

    return response.data;
}

export const useDeleteSchedule = () => {
    return useMutation({
        mutationFn: (scheduleId: string) => deleteSchedule(scheduleId),
        onSuccess: async() => {
            await new Promise(resolve => setTimeout(resolve, 500));
            queryClient.invalidateQueries({ queryKey: ['schedules'] });
        }
    });
};
