import type { TaskStatus } from '@/types';
import { Badge } from '@/components/ui/badge';
import { STATUS_LABELS, STATUS_COLORS } from '@/lib/constants';

export function StatusBadge({ status }: { status: TaskStatus }) {
  return <Badge variant="outline" className={STATUS_COLORS[status]}>{STATUS_LABELS[status]}</Badge>;
}
