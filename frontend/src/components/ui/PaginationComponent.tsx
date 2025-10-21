import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationPrevious, PaginationNext, PaginationEllipsis } from './pagination'

interface PaginationProps {
    currentPage: number;
    setCurrentPage: (page: number) => void;
    totalPages: number;
}

const PaginationComponent = ({
    currentPage,
    setCurrentPage,
    totalPages,
}: PaginationProps) => {

    if (totalPages <= 1) {
        return null;
    }

    console.log("PaginationComponent rendered with currentPage:", currentPage, "totalPages:", totalPages);

    const pages = [];
    if (currentPage == 1) {
        for (let i = 1; i <= Math.min(totalPages, 3); i++) {
            pages.push(i);
        }
    }
    else if (currentPage == totalPages) {
        for (let i = Math.max(1, totalPages - 2); i <= totalPages; i++) {
            pages.push(i);
        }
    }
    else {
        for (let i = Math.max(1, currentPage - 1); i <= Math.min(totalPages, currentPage + 1); i++) {
            pages.push(i);
        }
    }

    return (
        <Pagination>
            <PaginationContent>
                <PaginationItem>
                    <PaginationPrevious href="#" onClick={() => { setCurrentPage(currentPage - 1) }} aria-disabled={currentPage == 1} />
                </PaginationItem>

                {totalPages > 3 && currentPage > 2 && (
                    <PaginationItem>
                        <PaginationEllipsis />
                    </PaginationItem>
                )}

                {pages.map((page) => (
                    <PaginationItem key={page}>
                        <PaginationLink href="#" onClick={() => { setCurrentPage(page) }} isActive={currentPage == page}>
                            {page}
                        </PaginationLink>
                    </PaginationItem>
                ))}

                {totalPages > 3 && currentPage < totalPages - 1 && (
                    <PaginationItem>
                        <PaginationEllipsis />
                    </PaginationItem>
                )}

                <PaginationItem>
                    <PaginationNext href="#" onClick={() => { setCurrentPage(currentPage + 1) }} aria-disabled={currentPage == totalPages} />
                </PaginationItem>
            </PaginationContent>
        </Pagination>
    )
}

export default PaginationComponent;
