import type { ProjectReport as ProjectReportType } from '@/types';
import { formatMinutes } from '@/lib/utils';
import { Card, CardContent } from '@/components/ui/card';

interface ProjectReportProps {
  report: ProjectReportType;
}

function StatCard({ value, label, color }: { value: string; label: string; color?: string }) {
  return (
    <Card className="py-4 gap-0">
      <CardContent className="flex flex-col items-center">
        <p className={`text-2xl font-bold ${color ?? 'text-foreground'}`}>{value}</p>
        <p className="text-xs text-muted-foreground mt-1">{label}</p>
      </CardContent>
    </Card>
  );
}

export function ProjectReport({ report }: ProjectReportProps) {
  const pct = Math.min(report.completion_pct, 100);

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <StatCard value={String(report.total_tasks)} label="Total Tasks" />
        <StatCard value={String(report.completed_tasks)} label="Completed" color="text-green-600" />
        <StatCard value={`${report.completion_pct.toFixed(0)}%`} label="Progress" color="text-blue-600" />
        <StatCard value={formatMinutes(report.total_minutes)} label="Time Logged" />
      </div>
      <div className="w-full bg-secondary rounded-full h-2">
        <div
          className="bg-primary h-2 rounded-full transition-all duration-500"
          style={{ width: `${pct}%` }}
        />
      </div>
    </div>
  );
}
