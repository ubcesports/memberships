import { useState } from "react";

export function useResettablePagination(anchorKey: string, initialLimit: number) {
  const [limit, setLimit] = useState(initialLimit);
  const [pagination, setPagination] = useState({ offset: 0, anchor: "" });

  const offset = pagination.anchor === anchorKey ? pagination.offset : 0;

  const setOffset = (nextOffset: number) => {
    setPagination({ offset: nextOffset, anchor: anchorKey });
  };

  const changeLimit = (nextLimit: number) => {
    setLimit(nextLimit);
    setPagination({ offset: 0, anchor: anchorKey });
  };

  const resetOffset = () => {
    setPagination({ offset: 0, anchor: anchorKey });
  };

  return { offset, limit, setOffset, changeLimit, resetOffset };
}
