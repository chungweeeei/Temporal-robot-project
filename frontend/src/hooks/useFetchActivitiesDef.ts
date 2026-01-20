import axios from "axios";
import { useQuery } from "@tanstack/react-query";
import type { ActivityDefinition } from "../types/activities";

async function fetchActivitiesDef() {

    const response = await axios.get<Promise<ActivityDefinition[]>>(
        "http://localhost:3000/api/v1/activities",
        {
            headers: {
                "Content-Type": "application/json",
            }
        }
    );

    if (response.status !== 200) {
        throw new Error(`Failed to fetch activities definitions: ${response.statusText}`);
    }

    return response.data;
}

export function useFetchActivitiesDef() {
    return useQuery<ActivityDefinition[]>({
        queryKey: ["activities"],
        queryFn: fetchActivitiesDef,
    })
}