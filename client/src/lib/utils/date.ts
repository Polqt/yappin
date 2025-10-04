import { format, formatDistance } from 'date-fns';

export const formatDate = (date: string): string => {
  return format(new Date(date), 'PPP');
};

export const formatTimeAgo = (date: string): string => {
  return formatDistance(new Date(date), new Date(), { addSuffix: true });
};
