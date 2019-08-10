<template>
    <v-container>
        <v-form>
            <div class="title">General</div>
            <v-container>
                <v-layout>
                    <v-flex xs11>
                        <v-text-field
                            v-model="config.name"
                            label="Configuration name"
                            placeholder="Organize pictures"
                            clearable
                            required
                            :error-messages="nameError"
                            :disabled="isSubmitting"
                        ></v-text-field>
                    </v-flex>
                </v-layout>
            </v-container>

            <div class="title">Source directory</div>
            <v-container>
                <v-layout align-center>
                    <v-flex xs11>
                        <v-text-field
                            v-model="config.src.dir"
                            label="Source directory path"
                            hint="Full path of the directory to organize"
                            placeholder="C:\Users\John\Pictures"
                            clearable
                            required
                            :error-messages="srcDirError"
                            :disabled="isSubmitting"
                        ></v-text-field>
                    </v-flex>
                    <v-flex xs>
                        <v-tooltip bottom>
                            <template v-slot:activator="{ on }">
                                <v-btn
                                    type="button"
                                    icon
                                    color="info"
                                    v-on="on"
                                    @click="selectSrcDir"
                                    :disabled="isSubmitting"
                                >
                                    <v-icon medium>folder_open</v-icon>
                                </v-btn>
                            </template>
                            <span>Browse directories</span>
                        </v-tooltip>
                    </v-flex>
                </v-layout>
                <v-layout>
                    <v-flex xs11>
                        <v-checkbox
                            v-model="config.src.includeSubdirs"
                            label="Include subdirectories"
                            :disabled="isSubmitting"
                        ></v-checkbox>
                    </v-flex>
                </v-layout>
                <v-layout>
                    <v-flex xs11>
                        <v-select
                            v-model="config.src.defaultOpType"
                            :items="defaultOpTypeItems"
                            :disabled="isSubmitting"
                        ></v-select>
                    </v-flex>
                </v-layout>
            </v-container>

            <div class="title">Destination directories</div>
            <v-container>
                <v-layout
                    align-center
                    v-for="(dstDir, index) in config.dst.dirs"
                    :key="index"
                >
                    <v-flex xs2>
                        <v-text-field
                            :key="`config.dst.dirs.${index}.hotkey`"
                            v-model="dstDir.hotkey"
                            label="Hotkey"
                            placeholder="a"
                            clearable
                            maxLength="1"
                            required
                            :error-messages="hotkeyError(index)"
                            :disabled="isSubmitting"
                        ></v-text-field>
                    </v-flex>
                    <v-flex
                        v-bind="{
                            [`xs${config.dst.dirs.length == 1 ? 9 : 8}`]: true,
                        }"
                    >
                        <v-text-field
                            :key="`config.dst.dirs.${index}.dir`"
                            v-model="dstDir.dir"
                            label="Destination directory"
                            hint="Full path of a destination directory"
                            placeholder="C:\Users\John\CatPictures"
                            clearable
                            required
                            :error-messages="dstDirError(index)"
                            :disabled="isSubmitting"
                        ></v-text-field>
                    </v-flex>
                    <v-flex xs>
                        <v-tooltip bottom>
                            <template v-slot:activator="{ on }">
                                <v-btn
                                    type="button"
                                    icon
                                    color="info"
                                    v-on="on"
                                    @click="selectDstDir(index)"
                                    :disabled="isSubmitting"
                                >
                                    <v-icon medium>folder_open</v-icon>
                                </v-btn>
                            </template>
                            <span>Browse directories</span>
                        </v-tooltip>
                    </v-flex>
                    <v-flex xs v-if="config.dst.dirs.length > 1">
                        <v-tooltip bottom>
                            <template v-slot:activator="{ on }">
                                <v-btn
                                    type="button"
                                    icon
                                    color="error"
                                    v-on="on"
                                    @click="removeDstDir(index)"
                                    :disabled="isSubmitting"
                                >
                                    <v-icon medium>delete</v-icon>
                                </v-btn>
                            </template>
                            <span>Remove</span>
                        </v-tooltip>
                    </v-flex>
                </v-layout>
                <v-layout row>
                    <v-btn
                        type="button"
                        color="info"
                        @click="addDstDir"
                        :disabled="isSubmitting"
                    >
                        Add destination directory
                    </v-btn>
                </v-layout>
            </v-container>

            <div class="title">Advanced options</div>
            <v-container>
                <v-layout>
                    <v-flex xs11>
                        <v-text-field
                            v-model.number="config.ops.numWorkers"
                            label="Number of workers"
                            hint="Number of workers executing file operations concurrently"
                            type="number"
                            min="1"
                            max="5"
                            required
                            :error-messages="numWorkersError"
                            :disabled="isSubmitting"
                        ></v-text-field>
                    </v-flex>
                </v-layout>
                <v-layout>
                    <v-flex xs11>
                        <v-text-field
                            v-model.number="config.ops.maxTries"
                            label="Maximum file operation retries"
                            hint="Number of maximum  retries for a failed file operation"
                            type="number"
                            min="1"
                            max="1000000"
                            required
                            :error-messages="maxTriesError"
                            :disabled="isSubmitting"
                        ></v-text-field>
                    </v-flex>
                </v-layout>
            </v-container>

            <v-btn
                right
                color="success"
                large
                @click="submit"
                :loading="isSubmitting"
            >
                Create configuration
            </v-btn>
        </v-form>
    </v-container>
</template>


<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import router from '@/router';
import to from 'await-to-js';

import { Config, ConfigValidationError } from '@/api/config';
import { OpType } from '@/api/operation';
import { configValidatorAPI, dialogAPI } from '@/api/api';
import { capitalize } from '@/utils/utils';
import { organizer } from '@/store/modules/organizer';

const defaultConfigName = 'Anonymous configuration';
const defaultNumWorkers = 1;
const defaultMaxTries = 10000;

@Component
export default class ConfigForm extends Vue {
    @organizer.Action
    public loadConfig!: (config: Config) => Promise<void>;

    public isSubmitting: boolean = false;

    public validationErrors: ConfigValidationError = { errors: {} };

    public config: Config = {
        id: 0,
        name: '',
        src: {
            dir: '',
            includeSubdirs: false,
            defaultOpType: OpType.Copy,
        },
        dst: {
            dirs: [{ hotkey: '', dir: '' }],
        },
        ops: {
            numWorkers: defaultNumWorkers,
            maxTries: defaultMaxTries,
        },
    };

    public defaultOpTypeItems = [
        {
            text: 'Copy source files to destination directories',
            value: OpType.Copy,
        },
        {
            text: 'Move source files to destination directories',
            value: OpType.Move,
        },
    ];

    public async selectSrcDir() {
        const [_, dir] = await to<string, string>(dialogAPI.selectDirectory());
        if (dir) {
            this.config.src.dir = dir;
        }
    }

    public async selectDstDir(index: number) {
        const [_, dir] = await to<string, string>(dialogAPI.selectDirectory());
        if (dir) {
            this.config.dst.dirs[index].dir = dir;
        }
    }

    public addDstDir() {
        this.config.dst.dirs.push({ hotkey: '', dir: '' });
    }

    public removeDstDir(index: number) {
        const remaining = this.config.dst.dirs.filter((_, i) => i !== index);
        this.config.dst.dirs = remaining;
    }

    public async submit() {
        this.isSubmitting = true;

        // Substitute possibly wrong values.
        this.config.name = this.config.name || defaultConfigName;
        this.config.ops.numWorkers =
            this.config.ops.numWorkers || defaultNumWorkers;
        this.config.ops.maxTries = this.config.ops.maxTries || defaultMaxTries;

        const [err, _] = await to<void, string>(
            configValidatorAPI.validateConfig(this.config),
        );
        if (err) {
            const errors: ConfigValidationError = JSON.parse(err);
            this.validationErrors.errors = errors.errors;
            this.isSubmitting = false;
            return;
        }

        this.validationErrors.errors = {};
        await this.loadConfig(this.config);
        this.isSubmitting = false;
        router.push({ name: 'organize' });
    }

    public get nameError() {
        const err = this.validationErrors.errors['config.name'] || '';
        return capitalize(err);
    }

    public get srcDirError() {
        const err = this.validationErrors.errors['config.src.dir'] || '';
        return capitalize(err);
    }

    public hotkeyError(index: number) {
        const err =
            this.validationErrors.errors[`config.dst.dirs.${index}.hotkey`] ||
            '';
        return capitalize(err);
    }

    public dstDirError(index: number) {
        const err =
            this.validationErrors.errors[`config.dst.dirs.${index}.dir`] || '';
        return capitalize(err);
    }

    public get numWorkersError() {
        const err = this.validationErrors.errors['config.ops.numWorkers'] || '';
        return capitalize(err);
    }

    public get maxTriesError() {
        const err = this.validationErrors.errors['config.ops.maxTries'] || '';
        return capitalize(err);
    }
}
</script>
