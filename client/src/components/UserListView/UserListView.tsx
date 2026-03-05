import { useState } from 'react';

import Frame from '../common/Frame/Frame';
import GenericTable from '../common/GenericTable/GenericTable';


type Line = {
    id: number;
    name: string;
    size: number;
};

/* ----- Sample data ----- */
const INITIAL_LINES: Line[] = [
    { id: 3, name: 'Alpha', size: 42 },
    { id: 5, name: 'Bravo', size: 27 },
    { id: 1, name: 'Charlie', size: 35 },
    { id: 2, name: 'Delta', size: 19 },
    { id: 4, name: 'Echo', size: 50 },
];

function UserListView() {
    const [data] = useState<Line[]>(INITIAL_LINES);

    const handleRowClick = (line: Line) => {
        console.log('Clicked line name:', line.name);
    };

    return (
        <Frame>
            <div className="table-container">
                <h2>All Users</h2>
                <GenericTable
                    data={data}
                    columns={['name', 'id', 'size']}
                    onRowClick={handleRowClick}
                    filterable={true}
                />
            </div>
        </Frame>
    );
};

export default UserListView;
