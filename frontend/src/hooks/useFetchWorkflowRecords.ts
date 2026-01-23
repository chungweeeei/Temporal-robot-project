import axios from "axios";
import { useQuery } from "@tanstack/react-query";
import type { WorkflowRecord } from "@/types/workflows";


async function fetchWorkflowRecords(): Promise<WorkflowRecord[]> {

    const response = await axios.get<WorkflowRecord[]>(
        `http://localhost:3000/api/v1/workflows/records`,
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    )

    if (response.status !== 200){
        throw new Error(`Failed to fetch workflow status: ${response.statusText}`);
    }
    
    return response.data;
}


export const useFetchWorkflowRecords = () => {
    return useQuery({
        queryKey: ['workflows', "records"],
        queryFn: () => fetchWorkflowRecords(),
        refetchInterval: 5000,
    });
}
