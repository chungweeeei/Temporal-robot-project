import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Calendar, Clock } from "lucide-react";

export function ScheduleSection() {
  return (
    <Card className="mt-8">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Calendar className="h-5 w-5" />
          Scheduled Workflows
        </CardTitle>
        <CardDescription>
          設定自動排程執行 Workflow
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col items-center justify-center py-12 text-center">
          <div className="rounded-full bg-muted p-4 mb-4">
            <Clock className="h-8 w-8 text-muted-foreground" />
          </div>
          <h3 className="font-semibold text-lg mb-2">Coming Soon</h3>
          <p className="text-muted-foreground text-sm max-w-md">
            排程功能即將推出。你將可以設定 Workflow 在特定時間或週期性自動執行。
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
