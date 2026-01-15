import { useQuery } from "@tanstack/react-query";
import { fetchWorkflows } from "../utils/http";

export const useFetchWorkflow = () => {
    return useQuery({
        queryKey: ['workflows'],
        queryFn: fetchWorkflows,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });
}

