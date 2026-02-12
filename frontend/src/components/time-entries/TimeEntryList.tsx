import type { TimeEntry } from '@/types';
import { formatDate, formatMinutes } from '@/lib/utils';
import { TimeEntryForm } from './TimeEntryForm';
import { Button } from '@/components/ui/button';

interface TimeEntryListProps {
  entries: TimeEntry[];
  onAdd: (data: { minutes: number; description?: string; date: string }) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
}

export function TimeEntryList({ entries, onAdd, onDelete }: TimeEntryListProps) {
  const total = entries.reduce((sum, e) => sum + e.minutes, 0);

  return (
    <div className="space-y-4">
      <TimeEntryForm onSubmit={onAdd} />
      {entries.length === 0 ? (
        <p className="text-sm text-gray-400 text-center py-4">No time logged yet</p>
      ) : (
        <>
          <p className="text-sm text-gray-600 font-medium">Total: {formatMinutes(total)}</p>
          <div className="space-y-2">
            {entries.map((e) => (
              <div key={e.id} className="flex items-center justify-between bg-gray-50 rounded-md px-3 py-2">
                <div>
                  <span className="text-sm font-medium text-gray-900">{formatMinutes(e.minutes)}</span>
                  {e.description && (
                    <span className="text-sm text-gray-500 ml-2">{e.description}</span>
                  )}
                  <span className="text-xs text-gray-400 ml-2">{formatDate(e.date)}</span>
                </div>
                <Button variant="ghost" className="text-xs text-red-500" onClick={() => onDelete(e.id)}>
                  Delete
                </Button>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
