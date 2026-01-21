import { memo } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Calendar, PauseCircle, PlayCircle, Clock, Timer, Loader2, Trash2 } from "lucide-react";
import { CreateScheduleModal } from "./CreateScheduleModal";
import { useFetchSchedules } from "@/hooks/useFetchSchedules";
import { usePauseSchedule, useUnpauseSchedule } from "@/hooks/useControlSchedule";
import { useDeleteSchedule } from "@/hooks/useDeleteSchedule";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import type { Schedule, ScheduleRange } from "@/types/schema";

const WEEKDAYS = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

function getDayName(day: number) {
  return WEEKDAYS[day] || day.toString();
}

function getMonthName(month: number) {
  return MONTHS[month - 1] || month.toString();
}

function describeRange(
  ranges: ScheduleRange[] | undefined,
  min: number,
  max: number,
  unitSingular: string,
  unitPlural: string,
  valueMap?: (v: number) => string
): string | null {
  if (!ranges || ranges.length === 0) return null;

  if (ranges.length === 1) {
    const r = ranges[0];
    if (r.start === min && r.end === max && r.step > 1) {
      return `Every ${r.step} ${unitPlural}`;
    }
    if (r.start === min && r.end === max && r.step === 1) {
      return null; 
    }
  }

  const parts = ranges.map(r => {
    const val = (v: number) => (valueMap ? valueMap(v) : v.toString());
    if (r.start === r.end) return val(r.start);
    if (r.step === 1) return `${val(r.start)}-${val(r.end)}`;
    return `${val(r.start)}-${val(r.end)}/${r.step}`;
  });

  switch (unitSingular) {
    case "minute":
      return `at min ${parts.join(", ")}`;
    case "hour":
      return `at hour ${parts.join(", ")}`;
    case "second":
      return null; 
    case "day of week":
      return `on ${parts.join(", ")}`;
    case "month":
      return `in ${parts.join(", ")}`;
    case "day of month":
      return `on day ${parts.join(", ")}`;
    case "year":
      return `in year ${parts.join(", ")}`;
    default:
      return `${unitSingular} ${parts.join(", ")}`;
  }
}

function formatSpec(spec?: Schedule['spec']) {
  if (!spec) return 'No schedule details';
  
  if (spec.cron_expressions && spec.cron_expressions.length > 0) {
    return `Cron: ${spec.cron_expressions.join(', ')}`;
  }
  
  if (spec.calendars && spec.calendars.length > 0) {
    return spec.calendars.map(cal => {
      if (cal.comment) return cal.comment;
      
      const parts: string[] = [];

      const minDesc = describeRange(cal.minute, 0, 59, "minute", "minutes");
      const hourDesc = describeRange(cal.hour, 0, 23, "hour", "hours");
      
      if (minDesc) parts.push(minDesc);
      if (hourDesc) parts.push(hourDesc);

      const wDayDesc = describeRange(cal.day_of_week, 0, 6, "day of week", "days", getDayName);
      const mDayDesc = describeRange(cal.day_of_month, 1, 31, "day of month", "days");
      const monthDesc = describeRange(cal.month, 1, 12, "month", "months", getMonthName);
      
      if (wDayDesc) parts.push(wDayDesc);
      if (mDayDesc) parts.push(mDayDesc);
      if (monthDesc) parts.push(monthDesc);

      if (parts.length === 0) return "Runs continuously (Default)";
      
      return parts.join(", ");
    }).join('; ');
  }
  return 'No schedule details';
}

function formatDate(dateStr?: string) {
  if (!dateStr) return null;
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return null;
  return new Intl.DateTimeFormat('en-US', {
    month: 'short', day: 'numeric',
    hour: 'numeric', minute: 'numeric', second: 'numeric', hour12: false
  }).format(date);
}

// 獨立的 ScheduleCard component
interface ScheduleCardProps {
  schedule: Schedule;
}

const ScheduleCard = memo(function ScheduleCard({ schedule }: ScheduleCardProps) {

  const pauseSchedule = usePauseSchedule();
  const unpauseSchedule = useUnpauseSchedule();
  const deleteSchedule = useDeleteSchedule();

  const isPending = pauseSchedule.isPending || unpauseSchedule.isPending;

  const handleTogglePause = () => {
    if (schedule.paused) {
      unpauseSchedule.mutate(schedule.schedule_id);
    } else {
      pauseSchedule.mutate(schedule.schedule_id);
    }
  };

  const handleDelete = () => {
    if (confirm("Are you sure you want to delete this schedule?")) {
      deleteSchedule.mutate(schedule.schedule_id);
    }
  };

  return (
    <div className="flex items-center justify-between p-4 border rounded-lg bg-card">
      <div className="space-y-1">
        <div className="font-medium flex items-center gap-2">
          {schedule.schedule_id}
          <Badge variant={schedule.paused ? "secondary" : "default"} className="text-xs">
            {schedule.paused ? "Paused" : "Running"}
          </Badge>
        </div>
        <div className="text-sm text-muted-foreground flex items-center gap-2">
          <span className="font-semibold">Schedule:</span>
          <code className="bg-muted px-1.5 py-0.5 rounded text-xs font-mono">
            {formatSpec(schedule.spec)}
          </code>
        </div>

        <div className="flex flex-col gap-1 mt-2 sm:flex-row sm:gap-4 sm:items-center">
          {schedule.recent_run && (
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground/80" title="Last Execute Time">
              <Clock className="w-3.5 h-3.5" />
              <span>Last: {formatDate(schedule.recent_run)}</span>
            </div>
          )}
          {schedule.upcoming_run && (
            <div className="flex items-center gap-1.5 text-xs text-blue-600/80 dark:text-blue-400" title="Next Scheduled Run">
              <Timer className="w-3.5 h-3.5" />
              <span>Next: {formatDate(schedule.upcoming_run)}</span>
            </div>
          )}
        </div>
      </div>
      <div className="flex items-center gap-2">
        <Button
          variant="ghost"
          size="icon"
          onClick={handleTogglePause}
          disabled={isPending}
          className="cursor-pointer hover:bg-muted"
          title={schedule.paused ? "Resume Schedule" : "Pause Schedule"}
        >
          {isPending ? (
            <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
          ) : schedule.paused ? (
            <PlayCircle className="h-5 w-5 text-green-500" />
          ) : (
            <PauseCircle className="h-5 w-5 text-muted-foreground" />
          )}
        </Button>
        <Button
          variant="ghost"
          size="icon"
          onClick={handleDelete}
          disabled={deleteSchedule.isPending}
          className="cursor-pointer hover:bg-muted text-red-500 hover:text-red-700 hover:bg-red-100"
          title="Delete Schedule"
        >
          {deleteSchedule.isPending ? (
            <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
          ) : (
            <Trash2 className="h-5 w-5" />
          )}
        </Button>
      </div>
    </div>
  );
});

export function ScheduleSection() {
  const { data: schedules, isLoading, isError } = useFetchSchedules();

  return (
    <Card className="mt-8">
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle className="flex items-center gap-2">
            <Calendar className="h-5 w-5" />
            Scheduled Workflows
          </CardTitle>
          <CardDescription>
            setting up automated triggers for your workflows
          </CardDescription>
        </div>
        <CreateScheduleModal />
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="text-center py-8 text-muted-foreground">Loading schedules...</div>
        ) : isError ? (
          <div className="text-center py-8 text-red-500">Failed to load schedules</div>
        ) : schedules?.length === 0 ? (
          <div className="text-center py-12 text-muted-foreground border border-dashed rounded-lg">
            No schedules found. Create one nicely!
          </div>
        ) : (
          <div className="space-y-4">
            {schedules?.map((schedule) => (
              <ScheduleCard key={schedule.schedule_id} schedule={schedule} />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
