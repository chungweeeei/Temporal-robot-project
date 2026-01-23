import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { WorkflowStatusBadge } from "../shared/WorkflowStatusBadge";
import { Activity } from "lucide-react";
import { useFetchWorkflowRecords } from "@/hooks/useFetchWorkflowRecords";

export function ExecutionRecords() {
  
  const { data: records } = useFetchWorkflowRecords();
  
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center gap-2">
          <Activity className="h-5 w-5 text-primary" />
          <CardTitle>Recent Workflow Records</CardTitle>
        </div>
        <CardDescription>Latest workflow execution status</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="rounded-md border overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="bg-gray-50 text-gray-500 font-medium">
                <tr>
                  <th className="h-10 px-4 align-middle">Workflow ID</th>
                  <th className="h-10 px-4 align-middle">Run ID</th>
                  <th className="h-10 px-4 align-middle">Status</th>
                  <th className="h-10 px-4 align-middle">Start Time</th>
                  <th className="h-10 px-4 align-middle">End Time</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {!records?.length ? (
                   <tr>
                     <td colSpan={5} className="p-4 text-center text-gray-500">No records found</td>
                   </tr>
                ) : (
                  records.map((record) => (
                    <tr key={record.run_id} className="hover:bg-gray-50/50 transition-colors">
                      <td className="p-4 font-mono text-xs text-gray-600 max-w-[200px] truncate" title={record.workflow_id}>
                        {record.workflow_id}
                      </td>
                      <td className="p-4 font-mono text-xs text-gray-600" title={record.run_id}>
                        {record.run_id.slice(0, 8)}...
                      </td>
                      <td className="p-4">
                        <WorkflowStatusBadge status={record.status} />
                      </td>
                      <td className="p-4 text-gray-600 whitespace-nowrap">
                        {new Date(record.start_time).toLocaleString()}
                      </td>
                      <td className="p-4 text-gray-600 whitespace-nowrap">
                        {record.end_time ? new Date(record.end_time).toLocaleString() : '-'}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}