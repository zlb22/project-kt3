import { onMounted, ref } from 'vue'
import { defineStore } from 'pinia'
import { getConfig } from '@/api'


export const useFileStore = defineStore('file', () => {
    const audioFile = ref<File | null>(null);
    const imgFile = ref<File | null>(null);
    const commitStatus = ref(false)
    const uid = ref(0)
    const list = ref({})
    const confirmCommitStatus = ref(false)
    const items = [
        { number: 1, content: '方案绘制', selected: true },
        { number: 2, content: '思路阐述', selected: false },
        { number: 3, content: '完成', selected: false }
    ]

    const authToken = ref('')

    interface AxiosResponse {
        data: {
            school: School[],
            grade: Grade[]
        }
    }
    interface School {
        id: number;
        name: string;
    }
    interface Grade {
        id: number,
        name: string
    }
    onMounted(async () => {
        const { data } = await getConfig()
        list.value = data.data
    })

    function changeItemStep(index: number) {
        for (let i = 0; i < items.length; i++) {
            if (items[i].number <= index) {
                items[i].selected = true
            } else {
                items[i].selected = false
            }
        }
    }
    function setAudioFile(file: File) {
        audioFile.value = file;
        console.log(audioFile.value);
    }
    function setImgFile(file: File) {
        imgFile.value = file;
        console.log(imgFile.value);

    }
    function setCommitStatus(status: boolean) {
        commitStatus.value = status;
    }
    function setConfirmCommitStatus(status: boolean) {
        confirmCommitStatus.value = status;
    }
    function setUid(uuid: number) {
        uid.value = uuid;
    }
    function setAuthToken(token:string){
        authToken.value = token;
    }
    return {
        uid,
        list,
        audioFile,
        imgFile,
        commitStatus,
        items,
        confirmCommitStatus,
        authToken,
        setAuthToken,
        setUid,
        setAudioFile,
        setImgFile,
        setCommitStatus,
        changeItemStep,
        setConfirmCommitStatus
    }
})