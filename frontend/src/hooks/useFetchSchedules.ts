import { useQuery } from "@tanstack/react-query";
import { fetchSchedules } from "@/utils/http";
import type { Schedule } from "@/types/schema";

export function useFetchSchedules() {
  return useQuery<Schedule[]>({
    queryKey: ["schedules"],
    queryFn: fetchSchedules,
  });
}
