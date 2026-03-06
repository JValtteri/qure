import { useEffect, useState } from 'react';
import { signal } from '@preact/signals-react';

import Frame from '../common/Frame/Frame';
import GenericTable from '../common/GenericTable/GenericTable';
import { listAllClients, type ClientResponse } from '../../api/api';


const loadingClientList = signal(false);

interface Props {

}

function UserListView() {
    const [data, setData] = useState(new Array<ClientResponse>());

    const updateUserListHandler = updateUserList(setData);

    useEffect(() => {
        updateUserListHandler();
    }, []);

    const handleRowClick = (line: ClientResponse) => {
        console.log('Clicked line name:', line.Id);
    };

    return (
        <Frame>
            <div className="table-container">
                <h2>All Users</h2>
                <GenericTable
                    data={data}
                    columns={['Email', 'Role', 'IsTemporary', 'CreatedDt']}
                    onRowClick={handleRowClick}
                    filterable={true}
                    sortable={true}
                    interpretBigNumbersAs='date'
                />
            </div>
        </Frame>
    );
};

function updateUserList(setData: React.Dispatch<React.SetStateAction<ClientResponse[]>>): () => Promise<void> {
    return async () => {
        if (loadingClientList.value == true) {
            return;
        }
        loadingClientList.value = true;
        try {
            await listAllClients()
                .then(value => {
                    if (value != null) {
                        setData(value);
                    }
                });
        } catch (error: any) {
            console.warn(error.message);
        }
        loadingClientList.value = false;
    };
}

export default UserListView;
